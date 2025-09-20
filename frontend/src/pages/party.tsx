import { useParams } from "react-router";

export function Party() {
  const { partyId = '' } = useParams();
  return <>This is the party page for party {partyId}</>;
}