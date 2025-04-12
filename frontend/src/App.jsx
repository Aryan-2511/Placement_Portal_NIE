import {
  Navigate,
  RouterProvider,
  createBrowserRouter,
} from 'react-router-dom';
import { Toaster } from 'react-hot-toast';

// Layouts
import MainLayout from './components/layout/MainLayout'; // Unified layout for both student and admin
import GlobalStyles from './styles/GlobalStyles';

// Student Pages
import StudentDashboard from './pages/student/StudentDashboard';
import Opportunities from './pages/student/Opportunities';
import Applications from './pages/student/Applications';
import Profile from './pages/student/Profile';
import Feedback from './pages/student/Feedback';
import FAQ from './pages/student/FAQ';

// Admin Pages
import AdminDashboard from './pages/admin/AdminDashboard';
import ManageOpportunities from './pages/admin/ManageOpportunities';
import StudentDetails from './pages/admin/StudentDetails';

// Shared Pages
import NotFound from './pages/NotFound';
import ErrorPage from './pages/ErrorPage';

// Utilities
import PrivateRoute from './routes/PrivateRoute';
import LandingPage from './pages/LandingPage';
import OpportunityDetail from './features/opportunities/OpportunityDetail';
import { DarkModeProvider } from './context/DarkModeContext';
import { QueryClient } from '@tanstack/react-query';
import { QueryClientProvider } from '@tanstack/react-query';
import Schedule from './pages/admin/Schedule';
import Announcements from './pages/admin/Announcements';
import AdminPanel from './pages/admin/AdminPanel';
import AdminProfile from './pages/admin/AdminProfile';
import AddOpportunity from './features/opportunities/AddOpportunity';
import StudentProfile from './pages/admin/StudentProfile';
// import { UserRoleProvider } from './context/UserContext';
// import Logout from './features/authentication/Logout';

const router = createBrowserRouter([
  {
    path: '/',
    errorElement: <ErrorPage />, // Global error boundary
    children: [
      {
        index: true,
        element: <LandingPage />,
      },
      {
        element: <PrivateRoute />, // PrivateRoute applied at top level
        children: [
          {
            path: 'student',
            element: (
              <PrivateRoute role="STUDENT">
                <MainLayout />
              </PrivateRoute>
            ),
            children: [
              {
                index: true,
                element: <Navigate replace to="dashboard" />,
              },
              { path: 'dashboard', element: <StudentDashboard /> },
              { path: 'opportunities', element: <Opportunities /> },
              {
                path: 'opportunities/:opportunityId',
                element: <OpportunityDetail />,
              },
              { path: 'applications', element: <Applications /> },
              { path: 'profile', element: <Profile /> },
              { path: 'feedback', element: <Feedback /> },
              { path: 'faq', element: <FAQ /> },
            ],
          },
          {
            path: 'admin',
            element: (
              <PrivateRoute role="ADMIN">
                <MainLayout />
              </PrivateRoute>
            ),
            children: [
              { path: 'dashboard', element: <AdminDashboard /> },
              {
                path: 'manage_opportunities',
                element: <ManageOpportunities />,
              },
              {
                path: 'manage_opportunities/add_new_opportunity',
                element: <AddOpportunity />,
              },
              {
                path: 'manage_opportunities/:opportunityId',
                element: <OpportunityDetail />,
              },
              { path: 'profile', element: <AdminProfile /> },
              { path: 'student_details', element: <StudentDetails /> },
              { path: 'student_details/:usn', element: <StudentProfile /> },
              { path: 'schedule', element: <Schedule /> },
              { path: 'annoucements', element: <Announcements /> },
              { path: 'admin_panel', element: <AdminPanel /> },
            ],
          },
        ],
      },
      { path: '*', element: <NotFound /> },
    ],
  },
]);

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // staleTime: 60 * 1000,
      staleTime: 0,
    },
  },
});

const App = () => {
  return (
    <>
      <DarkModeProvider>
        {/* <UserRoleProvider> */}
        <QueryClientProvider client={queryClient}>
          <GlobalStyles />
          <Toaster
            position="top-center"
            reverseOrder={false}
            gutter={8}
            containerClassName=""
            containerStyle={{ margin: '8px' }}
            toastOptions={{
              // Define default options
              className: '',
              duration: 5000,
              style: {
                fontSize: '16px',
                maxWidth: '500px',
                padding: '16px 24px',
                backgroundColor: 'var(--color-grey-0)',
                color: 'var(--color-grey-700)',
              },

              // Default options for specific types
              success: {
                duration: 3000,
              },
              error: {
                duration: 5000,
              },
            }}
          />
          <RouterProvider router={router} />
        </QueryClientProvider>
        {/* </UserRoleProvider> */}
      </DarkModeProvider>
    </>
  );
};

export default App;
