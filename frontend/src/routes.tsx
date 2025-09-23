import { Group } from 'src/pages/group.page';
import { Home } from 'src/pages/home.page';

export const routes = [
  {
    path: '/',
    element: <Home />,
  },
  {
    path: '/groups/:groupId',
    element: <Group />,
  },
];
