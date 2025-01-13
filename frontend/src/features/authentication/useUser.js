import { useQueryClient } from '@tanstack/react-query';
import Cookies from 'js-cookie';

export const useUser = () => {
  const queryClient = useQueryClient();

  // Check if user data exists in the query cache
  const cachedUserData = queryClient.getQueryData(['user']);

  // If user data is not available in cache, fetch it from cookies
  if (!cachedUserData) {
    const userDataFromCookies = Cookies.get('user');

    if (userDataFromCookies) {
      const user = JSON.parse(userDataFromCookies);

      // Set the query cache with user data from cookies
      queryClient.setQueryData(['user'], user);
      return user;
    }
  } else {
    return cachedUserData;
  }
};
