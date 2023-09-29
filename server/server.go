package server

import (
	"context"
	"fmt"
	v1 "github.com/authzed/authzed-go/proto/authzed/api/v1"
	"github.com/authzed/authzed-go/v1"
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

func (p PrbacSpicedbServer) GetPrincipalAccess(ctx context.Context, request api.GetPrincipalAccessRequestObject) (api.GetPrincipalAccessResponseObject, error) {
	resp := api.GetPrincipalAccess200JSONResponse{}

	// TODO: the follow would orginarily be extracted from the request headers
	userOrg := "aspian"
	rootWorkspace := userOrg + "_root"

	servicePermissions := p.RbacServices[request.Params.Application]
	for key := range servicePermissions {
		permission := servicePermissions[key]

		// Step 1: If this user has checkpermission on root workspace, they get permission with no attribute filters

		r, err := p.SpicedbClient.CheckPermission(context.TODO(), &v1.CheckPermissionRequest{
			Resource: &v1.ObjectReference{
				ObjectType: "workspace",
				ObjectId:   rootWorkspace,
			},
			Permission: permission.Permission,
			Subject: &v1.SubjectReference{Object: &v1.ObjectReference{
				ObjectType: "user",
				ObjectId:   *request.Params.Username,
			}},
		})

		if err != nil {
			fmt.Errorf("Spicedb error: %v", err)
			return api.GetPrincipalAccess500JSONResponse{}, err
		}

		if r.Permissionship == v1.CheckPermissionResponse_PERMISSIONSHIP_HAS_PERMISSION {
			permTuple := request.Params.Application + ":" + key // of the form "playbook-dispatcher:run:read"

			resp.Data = append(resp.Data, api.Access{Permission: permTuple})

			continue // generic permission granted, so no need to look at attribute filters
		}

		// STEP 2: If they don't have unrestricted permission, check for attribute filtered permissions

	}

	return resp, nil
}

func (PrbacSpicedbServer) ListCrossAccountRequests(ctx context.Context, request api.ListCrossAccountRequestsRequestObject) (api.ListCrossAccountRequestsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) CreateCrossAccountRequests(ctx context.Context, request api.CreateCrossAccountRequestsRequestObject) (api.CreateCrossAccountRequestsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) GetCrossAccountRequest(ctx context.Context, request api.GetCrossAccountRequestRequestObject) (api.GetCrossAccountRequestResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) PatchCrossAccountRequest(ctx context.Context, request api.PatchCrossAccountRequestRequestObject) (api.PatchCrossAccountRequestResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) PutCrossAccountRequest(ctx context.Context, request api.PutCrossAccountRequestRequestObject) (api.PutCrossAccountRequestResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) ListGroups(ctx context.Context, request api.ListGroupsRequestObject) (api.ListGroupsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) CreateGroup(ctx context.Context, request api.CreateGroupRequestObject) (api.CreateGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) DeleteGroup(ctx context.Context, request api.DeleteGroupRequestObject) (api.DeleteGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) GetGroup(ctx context.Context, request api.GetGroupRequestObject) (api.GetGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) UpdateGroup(ctx context.Context, request api.UpdateGroupRequestObject) (api.UpdateGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) DeletePrincipalFromGroup(ctx context.Context, request api.DeletePrincipalFromGroupRequestObject) (api.DeletePrincipalFromGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) GetPrincipalsFromGroup(ctx context.Context, request api.GetPrincipalsFromGroupRequestObject) (api.GetPrincipalsFromGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) AddPrincipalToGroup(ctx context.Context, request api.AddPrincipalToGroupRequestObject) (api.AddPrincipalToGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) DeleteRoleFromGroup(ctx context.Context, request api.DeleteRoleFromGroupRequestObject) (api.DeleteRoleFromGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) ListRolesForGroup(ctx context.Context, request api.ListRolesForGroupRequestObject) (api.ListRolesForGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) AddRoleToGroup(ctx context.Context, request api.AddRoleToGroupRequestObject) (api.AddRoleToGroupResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) ListPermissions(ctx context.Context, request api.ListPermissionsRequestObject) (api.ListPermissionsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) ListPermissionOptions(ctx context.Context, request api.ListPermissionOptionsRequestObject) (api.ListPermissionOptionsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) ListPolicies(ctx context.Context, request api.ListPoliciesRequestObject) (api.ListPoliciesResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) CreatePolicies(ctx context.Context, request api.CreatePoliciesRequestObject) (api.CreatePoliciesResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) DeletePolicy(ctx context.Context, request api.DeletePolicyRequestObject) (api.DeletePolicyResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) GetPolicy(ctx context.Context, request api.GetPolicyRequestObject) (api.GetPolicyResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) UpdatePolicy(ctx context.Context, request api.UpdatePolicyRequestObject) (api.UpdatePolicyResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) ListPrincipals(ctx context.Context, request api.ListPrincipalsRequestObject) (api.ListPrincipalsResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) ListRoles(ctx context.Context, request api.ListRolesRequestObject) (api.ListRolesResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) CreateRole(ctx context.Context, request api.CreateRoleRequestObject) (api.CreateRoleResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) DeleteRole(ctx context.Context, request api.DeleteRoleRequestObject) (api.DeleteRoleResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) GetRole(ctx context.Context, request api.GetRoleRequestObject) (api.GetRoleResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) PatchRole(ctx context.Context, request api.PatchRoleRequestObject) (api.PatchRoleResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) UpdateRole(ctx context.Context, request api.UpdateRoleRequestObject) (api.UpdateRoleResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) GetRoleAccess(ctx context.Context, request api.GetRoleAccessRequestObject) (api.GetRoleAccessResponseObject, error) {
	//TODO implement me
	panic("implement me")
}

func (PrbacSpicedbServer) GetStatus(ctx context.Context, request api.GetStatusRequestObject) (api.GetStatusResponseObject, error) {
	//TODO implement me
	panic("implement me")
}
