import { Party } from "src/pages/party";
import { Home } from "src/pages/home";

export const routes = [
  {
    path: "/",
    element: <Home />
  },
  {
    path: "/:partyId",
    element: <Party />
  }
]