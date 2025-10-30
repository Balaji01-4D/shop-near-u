// Utility to update page title dynamically
export function updatePageTitle(title: string) {
  document.title = `${title} | ShopNearU`
}

export function setDefaultTitle() {
  document.title = 'ShopNearU - Discover Local Shops'
}

// Page titles for different routes
export const pageTitles = {
  home: 'Home',
  shops: 'Discover Shops',
  products: 'Products',
  subscriptions: 'My Subscriptions',
  login: 'Sign In',
  register: 'Create Account',
  shopDetails: 'Shop Details',
} as const