import { Authorizer } from 'casbin.js';
import { handleInlineError } from './errorHandling';

export interface ResourcePermissions {
  canCreate: boolean;
  canUpdate: boolean;
  canDelete: boolean;
  canRead: boolean;
}

export interface User {
  id: string;
}

export type PermissionAction = 'create' | 'view' | 'update' | 'delete';

/**
 * Checks permissions for a resource type in a given namespace
 */
export async function permissionChecker(user: User, resourceType: string, namespaceId: string, actions: PermissionAction[] = ['create', 'view', 'update', 'delete']): Promise<ResourcePermissions> {
  const permissions: ResourcePermissions = {
    canCreate: false,
    canRead: false,
    canUpdate: false,
    canDelete: false
  };

  try {
    const authorizer = new Authorizer('auto', {
      endpoint: '/api/v1/permissions'
    });
    await authorizer.setUser(`user:${user.id}`);

    // Check all requested permissions in parallel
    const results = await Promise.all(
      actions.map(action => authorizer.can(action, resourceType, namespaceId))
    );

    // Map results back to permissions
    actions.forEach((action, index) => {
      switch (action) {
        case 'create':
          permissions.canCreate = results[index];
          break;
        case 'view':
          permissions.canRead = results[index];
          break;
        case 'update':
          permissions.canUpdate = results[index];
          break;
        case 'delete':
          permissions.canDelete = results[index];
          break;
      }
    });
  } catch (err) {
    handleInlineError(err, 'Unable to Check Permissions');
    // permissions remain false on error
  }

  return permissions;
}