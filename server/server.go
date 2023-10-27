package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
	"github.com/google/uuid"
	"github.com/merlante/prbac-spicedb/api"
)

type Filter struct {
	Name         string `json:"name"`
	Operator     string `json:"operator"`
	ResourceType string `json:"resourceType"`
	Verb         string `json:"verb"`
}

type ResourcePerm struct {
	Permission string `json:"permission"`
	Filter     Filter `json:"filter"`
}

type Permission map[string]ResourcePerm

type Services map[string]Permission

type PrbacSpicedbServer struct {
	RbacServices  Services
	SpicedbClient *authzed.Client
}

var permissionsToSystemRoles = map[string]string{
	"inventory:hosts:read":   "e18257ae-7506-11ee-8c9d-0242ac170005",
	"inventory:hosts:write":  "c4eaa6fb-7506-11ee-8c9d-0242ac170005",
	"inventory:groups:read":  "10d23bda-7507-11ee-8c9d-0242ac170005",
	"inventory:groups:write": "fca60508-7506-11ee-8c9d-0242ac170005",
}

func (p *PrbacSpicedbServer) GetPrincipalAccess(ctx context.Context, request api.GetPrincipalAccessRequestObject) (api.GetPrincipalAccessResponseObject, error) {
	resp := api.GetPrincipalAccess200JSONResponse{}

	// TODO: the following would ordinarily be extracted from the request headers
	userOrg := "aspian"
	rootWorkspace := userOrg + "_root"

	servicePermissions := p.RbacServices[request.Params.Application]
	for key := range servicePermissions {
		servicePermission := servicePermissions[key]

		// Step 1: If this user has checkpermission on root workspace, they get permission with no attribute filters

		r, err := p.SpicedbClient.CheckPermission(ctx, &v1.CheckPermissionRequest{
			Resource: &v1.ObjectReference{
				ObjectType: "workspace",
				ObjectId:   rootWorkspace,
			},
			Permission: servicePermission.Permission,
			Subject: &v1.SubjectReference{Object: &v1.ObjectReference{
				ObjectType: "user",
				ObjectId:   *request.Params.Username,
			}},
		})

		if err != nil {
			fmt.Errorf("spicedb error: %v", err)
			return api.GetPrincipalAccess500JSONResponse{}, err
		}

		if r.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
			permTuple := request.Params.Application + ":" + key // of the form "playbook-dispatcher:run:read"

			resp.Data = append(resp.Data, api.Access{Permission: permTuple})

			continue // generic permission granted, so no need to look at attribute filters
		}

		// STEP 2: If they don't have unrestricted permission, check for attribute filtered permissions

		boundResourceType := servicePermission.Filter.ResourceType
		permission := servicePermission.Filter.Verb

		lrClient, err := p.SpicedbClient.LookupResources(ctx, &v1.LookupResourcesRequest{
			ResourceObjectType: boundResourceType,
			Permission:         permission,
			Subject: &v1.SubjectReference{
				Object: &v1.ObjectReference{
					ObjectType: "user",
					ObjectId:   *request.Params.Username,
				},
			},
		})

		if err != nil {
			fmt.Errorf("spicedb error: %v", err)
			return api.GetPrincipalAccess500JSONResponse{}, err
		}

		var boundResources []string
		for {
			next, err := lrClient.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				fmt.Errorf("spicedb error: %v", err)
				return api.GetPrincipalAccess500JSONResponse{}, err
			}

			boundResources = append(boundResources, next.GetResourceObjectId()) // e.g. service or inventory group
		}

		// STEP 2a: see if the user permission can be matched to the scope
		if len(boundResources) == 0 {
			continue
		}

		var resourceDefinitions []api.ResourceDefinition

		for _, resource := range boundResources {
			operator := servicePermission.Filter.Operator

			if strings.EqualFold("equal", operator) {
				filter := api.ResourceDefinitionFilter{
					Key:       servicePermission.Filter.Name,
					Operation: api.ResourceDefinitionFilterOperation(operator),
					Value:     resource,
				}

				resourceDefinitions = append(resourceDefinitions, api.ResourceDefinition{AttributeFilter: filter})

			} else {
				fmt.Errorf("unsupported PRBAC operator: %s", operator)
				continue
			}
		}

		if len(resourceDefinitions) != 0 {
			// We can make an attribute filter for each containing thing (service, inventory group) the user has access to
			permTuple := request.Params.Application + ":" + key // of the form "playbook-dispatcher:run:read"

			resp.Data = append(resp.Data, api.Access{
				Permission:          permTuple,
				ResourceDefinitions: resourceDefinitions,
			})
		}
	}

	return resp, nil
}

func (*PrbacSpicedbServer) ListCrossAccountRequests(ctx context.Context, request api.ListCrossAccountRequestsRequestObject) (api.ListCrossAccountRequestsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) CreateCrossAccountRequests(ctx context.Context, request api.CreateCrossAccountRequestsRequestObject) (api.CreateCrossAccountRequestsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) GetCrossAccountRequest(ctx context.Context, request api.GetCrossAccountRequestRequestObject) (api.GetCrossAccountRequestResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) PatchCrossAccountRequest(ctx context.Context, request api.PatchCrossAccountRequestRequestObject) (api.PatchCrossAccountRequestResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) PutCrossAccountRequest(ctx context.Context, request api.PutCrossAccountRequestRequestObject) (api.PutCrossAccountRequestResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) ListGroups(ctx context.Context, request api.ListGroupsRequestObject) (api.ListGroupsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) CreateGroup(ctx context.Context, request api.CreateGroupRequestObject) (api.CreateGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) DeleteGroup(ctx context.Context, request api.DeleteGroupRequestObject) (api.DeleteGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) GetGroup(ctx context.Context, request api.GetGroupRequestObject) (api.GetGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) UpdateGroup(ctx context.Context, request api.UpdateGroupRequestObject) (api.UpdateGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PrbacSpicedbServer) DeletePrincipalFromGroup(ctx context.Context, request api.DeletePrincipalFromGroupRequestObject) (api.DeletePrincipalFromGroupResponseObject, error) {
	userNames := strings.Split(request.Params.Usernames, ",")
	updates := make([]*v1.RelationshipUpdate, len(userNames))
	for i, username := range userNames {
		updates[i] = &v1.RelationshipUpdate{
			Operation: v1.RelationshipUpdate_OPERATION_DELETE,
			Relationship: &v1.Relationship{
				Resource: &v1.ObjectReference{
					ObjectType: "group",
					ObjectId:   request.Uuid.String(),
				},
				Relation: "member",
				Subject: &v1.SubjectReference{
					Object: &v1.ObjectReference{
						ObjectType: "user",
						ObjectId:   username, // TODO: needs to be an ID not a username
					},
				},
			},
		}
	}

	_, err := p.SpicedbClient.WriteRelationships(ctx, &v1.WriteRelationshipsRequest{
		Updates: updates,
	})

	if err != nil {
		return api.DeletePrincipalFromGroup500JSONResponse{}, err
	}

	return api.DeletePrincipalFromGroup204Response{}, nil
}

func (*PrbacSpicedbServer) GetPrincipalsFromGroup(ctx context.Context, request api.GetPrincipalsFromGroupRequestObject) (api.GetPrincipalsFromGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PrbacSpicedbServer) AddPrincipalToGroup(ctx context.Context, request api.AddPrincipalToGroupRequestObject) (api.AddPrincipalToGroupResponseObject, error) {
	updates := make([]*v1.RelationshipUpdate, len(request.Body.Principals))
	for i, principal := range request.Body.Principals {
		updates[i] = &v1.RelationshipUpdate{
			Operation: v1.RelationshipUpdate_OPERATION_TOUCH,
			Relationship: &v1.Relationship{

				Resource: &v1.ObjectReference{
					ObjectType: "group",
					ObjectId:   request.Uuid.String(),
				},
				Relation: "member",
				Subject: &v1.SubjectReference{
					Object: &v1.ObjectReference{
						ObjectType: "user",
						ObjectId:   principal.Username, // TODO: needs to be an ID not a username
					},
				},
			},
		}
	}

	_, err := p.SpicedbClient.WriteRelationships(ctx, &v1.WriteRelationshipsRequest{
		Updates: updates,
	})

	if err != nil {
		return api.AddPrincipalToGroup500JSONResponse{}, err
	}

	return api.AddPrincipalToGroup200JSONResponse{}, nil
}

func (*PrbacSpicedbServer) DeleteRoleFromGroup(ctx context.Context, request api.DeleteRoleFromGroupRequestObject) (api.DeleteRoleFromGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) ListRolesForGroup(ctx context.Context, request api.ListRolesForGroupRequestObject) (api.ListRolesForGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PrbacSpicedbServer) AddRoleToGroup(ctx context.Context, request api.AddRoleToGroupRequestObject) (api.AddRoleToGroupResponseObject, error) {
	updates := make([]*v1.RelationshipUpdate, 0)

	for _, role := range request.Body.Roles {
		client, err := p.SpicedbClient.LookupSubjects(ctx, &v1.LookupSubjectsRequest{
			Resource: &v1.ObjectReference{
				ObjectType: "rbac/v1role",
				ObjectId:   role.String(),
			},
			Permission:            "binding",
			SubjectObjectType:     "role_binding",
			OptionalConcreteLimit: 0,
		})

		if err != nil {
			return nil, err
		}

		for true {
			result, err := client.Recv()
			if err != nil {
				if err == io.EOF {
					break
				} else {
					return nil, err
				}
			}

			updates = append(updates, &v1.RelationshipUpdate{
				Operation: v1.RelationshipUpdate_OPERATION_TOUCH,
				Relationship: &v1.Relationship{
					Resource: &v1.ObjectReference{
						ObjectType: "role_binding",
						ObjectId:   result.Subject.SubjectObjectId,
					},
					Relation: "subject",
					Subject: &v1.SubjectReference{
						Object: &v1.ObjectReference{
							ObjectType: "group",
							ObjectId:   request.Uuid.String(),
						},
						OptionalRelation: "member",
					},
				},
			})
		}
	}

	//TODO: some sort of concurrency check is required here: this write is dependent upon results read above
	_, err := p.SpicedbClient.WriteRelationships(ctx, &v1.WriteRelationshipsRequest{
		Updates: updates,
	})

	if err != nil {
		return api.AddRoleToGroup500JSONResponse{}, err
	}

	return api.AddRoleToGroup200JSONResponse{}, nil
}

func (*PrbacSpicedbServer) ListPermissions(ctx context.Context, request api.ListPermissionsRequestObject) (api.ListPermissionsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) ListPermissionOptions(ctx context.Context, request api.ListPermissionOptionsRequestObject) (api.ListPermissionOptionsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) ListPolicies(ctx context.Context, request api.ListPoliciesRequestObject) (api.ListPoliciesResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) CreatePolicies(ctx context.Context, request api.CreatePoliciesRequestObject) (api.CreatePoliciesResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) DeletePolicy(ctx context.Context, request api.DeletePolicyRequestObject) (api.DeletePolicyResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) GetPolicy(ctx context.Context, request api.GetPolicyRequestObject) (api.GetPolicyResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) UpdatePolicy(ctx context.Context, request api.UpdatePolicyRequestObject) (api.UpdatePolicyResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) ListPrincipals(ctx context.Context, request api.ListPrincipalsRequestObject) (api.ListPrincipalsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) ListRoles(ctx context.Context, request api.ListRolesRequestObject) (api.ListRolesResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PrbacSpicedbServer) CreateRole(ctx context.Context, request api.CreateRoleRequestObject) (api.CreateRoleResponseObject, error) {
	// TODO: the following would ordinarily be extracted from the request headers
	userOrg := "aspian"
	rootWorkspace := userOrg + "_root"

	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	roleId := id.String()

	// Create corresponding role and rolebinding
	updates := make([]*v1.RelationshipUpdate, 4)

	updates[0] = createRelationshipUpdate(v1.RelationshipUpdate_OPERATION_TOUCH, "role_binding", roleId, "granted", "role", roleId)
	updates[1] = createRelationshipUpdate(v1.RelationshipUpdate_OPERATION_TOUCH, "workspace", rootWorkspace, "user_grant", "role_binding", roleId)
	updates[2] = createRelationshipUpdate(v1.RelationshipUpdate_OPERATION_TOUCH, "rbac/v1role", roleId, "role", "role", roleId)
	updates[3] = createRelationshipUpdate(v1.RelationshipUpdate_OPERATION_TOUCH, "rbac/v1role", roleId, "binding", "role_binding", roleId)

	for _, access := range request.Body.Access {
		if access.ResourceDefinitions == nil {
			//Add converted role permissions
			convertedPermission := cleanNameForSchemaCompatibility(access.Permission)
			updates = append(updates, createRelationshipUpdate(v1.RelationshipUpdate_OPERATION_TOUCH, "role", roleId, convertedPermission, "user", "*"))
		} else {
			for _, definition := range access.ResourceDefinitions {
				filter := definition.AttributeFilter

				switch filter.Key {
				case "group.id":
					if role, ok := permissionsToSystemRoles[access.Permission]; ok {
						bindingId := roleId + "_" + filter.Value //TODO: value can be an array, but the generated API doesn't accept it.

						updates = append(updates, createRelationshipUpdate(v1.RelationshipUpdate_OPERATION_TOUCH, "workspace", filter.Value, "parent", "workspace", rootWorkspace))
						updates = append(updates, createRelationshipUpdate(v1.RelationshipUpdate_OPERATION_TOUCH, "workspace", filter.Value, "user_grant", "role_binding", bindingId))
						updates = append(updates, createRelationshipUpdate(v1.RelationshipUpdate_OPERATION_TOUCH, "rbac/v1role", roleId, "binding", "role_binding", bindingId))
						updates = append(updates, createRelationshipUpdate(v1.RelationshipUpdate_OPERATION_TOUCH, "role_binding", bindingId, "granted", "role", role))
					}
				default:
					fmt.Printf("[INFO] Unhandled resource definition for permission %s in role %s, key: %s\n", access.Permission, request.Body.Name, filter.Key)
				}
			}
		}
	}

	_, err = p.SpicedbClient.WriteRelationships(ctx, &v1.WriteRelationshipsRequest{
		Updates: updates,
	})

	if err != nil {
		return nil, err
	}

	return api.CreateRole201JSONResponse{
		Uuid: id,
	}, nil
}

func (*PrbacSpicedbServer) DeleteRole(ctx context.Context, request api.DeleteRoleRequestObject) (api.DeleteRoleResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) GetRole(ctx context.Context, request api.GetRoleRequestObject) (api.GetRoleResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) PatchRole(ctx context.Context, request api.PatchRoleRequestObject) (api.PatchRoleResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (*PrbacSpicedbServer) UpdateRole(ctx context.Context, request api.UpdateRoleRequestObject) (api.UpdateRoleResponseObject, error) {
	//TODO implement me
	//Needs to perform similar writes to CreateRole ~plus~:
	//* Any no longer used permissions need to be revoked, see AddGroupToRole for how to find existing role_bindings
	//* Any new role_bindings that need to be created need their subject to be copied from the original role_binding subject
	panic("implement me")
}

func (p *PrbacSpicedbServer) GetRoleAccess(ctx context.Context, request api.GetRoleAccessRequestObject) (api.GetRoleAccessResponseObject, error) {
	// TODO: Not fully tested/implemented -- see discussion in getPRBACPermsFromSpicedbPerms

	resp := api.GetRoleAccess200JSONResponse{}

	role := request.Uuid // assume that the uuid form is the form that we are storing in spicedb

	rClient, err := p.SpicedbClient.ReadRelationships(ctx, &v1.ReadRelationshipsRequest{
		RelationshipFilter: &v1.RelationshipFilter{
			ResourceType:       "role",
			OptionalResourceId: role.String(),
			OptionalSubjectFilter: &v1.SubjectFilter{
				SubjectType: "user",
			},
		},
	})

	if err != nil {
		fmt.Errorf("spicedb error: %v", err)
		return api.GetRoleAccess500JSONResponse{}, err
	}

	var relationships []string
	for {
		next, err := rClient.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Errorf("spicedb error: %v", err)
			return api.GetRoleAccess500JSONResponse{}, err
		}

		relationships = append(relationships, next.GetRelationship().GetRelation())
	}

	resp.Data = p.getPRBACPermsFromSpicedbPerms(relationships)

	return resp, nil
}

func (*PrbacSpicedbServer) GetStatus(ctx context.Context, request api.GetStatusRequestObject) (api.GetStatusResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (p *PrbacSpicedbServer) getPRBACPermsFromSpicedbPerms(spicedbPerms []string) (accesses []api.Access) {
	var spiceToPRbacMapping map[string]string

	for service, permission := range p.RbacServices {
		for rbacPerm, resourcePerm := range permission {
			spiceDbPerm := resourcePerm.Permission

			spiceToPRbacMapping[spiceDbPerm] = service + ":" + rbacPerm
		}
	}

	for _, spiceDbPerm := range spicedbPerms {
		rbacPerm, mappingFound := spiceToPRbacMapping[spiceDbPerm]

		if mappingFound {
			accesses = append(accesses, api.Access{
				Permission: rbacPerm,
				// TODO: anything for resource definitions?

				// Discussion:
				// Without adding roles to the filter in the config (see below), it's not clear how we would be able to attach resourceDefinitions to a role.
				// However, adding roles may also imply changes in spicedb, bringing up a host of consistency issues we'd like to avoid.
				//
				//{
				//	"playbook-dispatcher": {
				//		"run:read": {
				//			"permission": "dispatcher_view_runs",
				//			"filter": {
				//				"roles": ["task_admin", "remediations_admin"],
				//				"name": "service",
				//				"operator": "equal",
				//				"resourceType": "dispatcher/service",
				//				"verb": "view"
				//			}
				//		}
				//	}
				//}

				// But before we do something like this, we want to step back and see if this endpoint enjoys enough usage to warrant it.
			})
		}
	}

	return
}

func cleanNameForSchemaCompatibility(name string) string { //Taken from schema translator
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, ":", "_")
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "*", "all")

	return name
}

func createRelationshipUpdate(operation v1.RelationshipUpdate_Operation, objectType, objectId, relation, subjectType, subjectId string) *v1.RelationshipUpdate {
	return &v1.RelationshipUpdate{
		Operation: v1.RelationshipUpdate_OPERATION_TOUCH,
		Relationship: &v1.Relationship{
			Resource: &v1.ObjectReference{
				ObjectType: objectType,
				ObjectId:   objectId,
			},
			Relation: relation,
			Subject: &v1.SubjectReference{
				Object: &v1.ObjectReference{
					ObjectType: subjectType,
					ObjectId:   subjectId,
				},
			},
		},
	}
}
