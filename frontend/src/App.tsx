import { createBrowserRouter } from 'react-router';
import { RouterProvider } from 'react-router/dom';

import { routes } from 'src/routes';

function App() {
  const router = createBrowserRouter(routes);

  return (
    <div className='App'>
      <RouterProvider router={router} />
    </div>
  );
}

export default App;
