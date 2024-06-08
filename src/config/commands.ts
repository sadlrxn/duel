export type UserRole = 'user' | 'ambassador' | 'moderator' | 'admin';

export interface DuelCommand {
  pattern: string;
  regExp: RegExp;
  role: UserRole;
}

export const convertUserPermissiontoNumber = (role: UserRole) => {
  switch (role) {
    case 'ambassador':
      return 1;
    case 'moderator':
      return 2;
    case 'admin':
      return 3;
    default:
      return 0;
  }
};

export const checkUserHasPermission = (
  userRole: UserRole,
  requiredRole: UserRole
) => {
  return (
    convertUserPermissiontoNumber(userRole) >=
    convertUserPermissiontoNumber(requiredRole)
  );
};
