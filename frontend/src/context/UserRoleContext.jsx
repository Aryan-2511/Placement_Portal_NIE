import { createContext, useContext } from 'react';

const UserRoleContext = createContext();

const UserRoleProvider = ({ children, role }) => {
  return (
    <UserRoleContext.Provider value={{ role }}>
      {children}
    </UserRoleContext.Provider>
  );
};
const useRole = () => {
  const context = useContext(UserRoleContext);
  if (context === undefined)
    throw new Error('UserRoleContext was used outside of UserRoleProvider');
  return context;
};
export { UserRoleProvider, useRole };
