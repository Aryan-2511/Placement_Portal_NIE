import { RouterProvider, createBrowserRouter } from 'react-router-dom';

// Layouts
import MainLayout from './components/layout/MainLayout'; // Unified layout for both student and coordinator

import GlobalStyles from './styles/GlobalStyles';

// Student Pages
import StudentDashboard from './pages/student/StudentDashboard';
import Opportunities from './pages/student/Opportunities';
import Applications from './pages/student/Applications';
import Profile from './pages/student/Profile';
import Feedback from './pages/student/Feedback';
import FAQ from './pages/student/FAQ';

// Coordinator Pages
import CoordinatorDashboard from './pages/coordinator/CoordinatorDashboard';
import ManageOpportunities from './pages/coordinator/ManageOpportunities';
import Reports from './pages/coordinator/Reports';

// Shared Pages
import Login from './pages/Login';
import NotFound from './pages/NotFound';
import ErrorPage from './pages/ErrorPage';

// Utilities
import PrivateRoute from './routes/PrivateRoute';

const router = createBrowserRouter([
  {
    path: '/',
    errorElement: <ErrorPage />, // Global error boundary
    children: [
      {
        index: true,
        element: <Login />,
      },
      {
        element: <PrivateRoute />, // PrivateRoute applied at top level
        children: [
          {
            path: 'student',
            element: (
              <PrivateRoute role="student">
                <MainLayout />
              </PrivateRoute>
            ),
            children: [
              { index: true, path: 'dashboard', element: <StudentDashboard /> },
              { path: 'opportunities', element: <Opportunities /> },
              { path: 'applications', element: <Applications /> },
              { path: 'profile', element: <Profile /> },
              { path: 'feedback', element: <Feedback /> },
              { path: 'faq', element: <FAQ /> },
            ],
          },
          {
            path: 'coordinator',
            element: (
              <PrivateRoute role="coordinator">
                <MainLayout />
              </PrivateRoute>
            ),
            children: [
              { path: 'dashboard', element: <CoordinatorDashboard /> },
              {
                path: 'manage-opportunities',
                element: <ManageOpportunities />,
              },
              { path: 'reports', element: <Reports /> },
            ],
          },
        ],
      },
      { path: '*', element: <NotFound /> },
    ],
  },
]);

const App = () => {
  return (
    <>
      <GlobalStyles />
      <RouterProvider router={router} />;
    </>
  );
};

export default App;
