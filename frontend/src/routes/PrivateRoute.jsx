// src/routes/PrivateRoute.jsx
// import { Navigate } from "react-router-dom";
// import { useSelector } from "react-redux";

import { Outlet } from 'react-router-dom';

// const PrivateRoute = ({ children, role }) => {
//   const { isAuthenticated, userRole } = useSelector((state) => state.auth);

//   if (!isAuthenticated) {
//     return <Navigate to="/" replace />;
//   }

//   if (role && userRole !== role) {
//     // Redirect to login or a "not authorized" page if the role does not match
//     return <Navigate to="/" replace />;
//   }

//   return children;
// };

// export default PrivateRoute;

function PrivateRoute({ children, role = 'role' }) {
  return <>{children || <Outlet />}</>;
}

export default PrivateRoute;
