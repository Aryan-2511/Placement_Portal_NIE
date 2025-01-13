import { createContext, useContext, useState } from 'react';

const UserContext = createContext();

const UserProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  return (
    <UserContext.Provider value={{ user, setUser }}>
      {children}
    </UserContext.Provider>
  );
};
const useUser = () => {
  const context = useContext(UserContext);
  if (!context)
    throw new Error('useUser must be used within a UserRoleProvider');
  return context;
};
export { UserProvider, useUser };
