import { Group } from "src/pages/group";
import { Home } from "src/pages/home";

export const routes = [
  {
    path: "/",
    element: <Home />
  },
  {
    path: "/:groupId",
    element: <Group />
  }
]