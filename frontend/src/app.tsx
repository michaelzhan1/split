import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

import { createBrowserRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';
import { routes } from 'src/routes';

function App() {
  const router = createBrowserRouter(routes);
  const queryClient = new QueryClient();

  return (
    <div className='app'>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </div>
  );
}

export default App;
