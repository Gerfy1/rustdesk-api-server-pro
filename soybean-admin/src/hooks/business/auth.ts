import { computed } from 'vue';
import { useAuthStore } from '@/store/modules/auth';

/**
 * User role levels
 */
export const USER_ROLES = {
  USER: 1,          // Regular user - RustDesk client only
  SUPPORT: 2,       // Support - View logs and devices
  SUPPORT_N2: 3,    // Support N2 - Manage users and address books
  SUPER_ADMIN: 4    // Super Admin - Full access
} as const;

export function useAuth() {
  const authStore = useAuthStore();

  function hasAuth(codes: string | string[]) {
    if (!authStore.isLogin) {
      return false;
    }

    if (typeof codes === 'string') {
      return authStore.userInfo.buttons.includes(codes);
    }

    return codes.some(code => authStore.userInfo.buttons.includes(code));
  }

  return {
    hasAuth
  };
}

/**
 * Check if user has required role level
 */
export function usePermission() {
  const authStore = useAuthStore();
  
  const userRole = computed(() => {
    const role = authStore.userInfo.role || 0;
    console.log('[Permission] User role:', role, 'UserInfo:', authStore.userInfo);
    return role;
  });
  
  const hasPermission = (requiredRole: number): boolean => {
    const hasAccess = userRole.value >= requiredRole;
    console.log('[Permission] Required:', requiredRole, 'User has:', userRole.value, 'Access:', hasAccess);
    return hasAccess;
  };
  
  const canManageAddressBooks = computed(() => hasPermission(USER_ROLES.SUPPORT_N2));
  const canManageUsers = computed(() => hasPermission(USER_ROLES.SUPPORT_N2));
  const canViewOnly = computed(() => userRole.value === USER_ROLES.SUPPORT);
  const isSuperAdmin = computed(() => userRole.value === USER_ROLES.SUPER_ADMIN);
  
  return {
    userRole,
    hasPermission,
    canManageAddressBooks,
    canManageUsers,
    canViewOnly,
    isSuperAdmin,
    USER_ROLES
  };
}
