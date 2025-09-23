import { useParams } from 'react-router';

export function Group() {
  const { groupId = '' } = useParams();
  return <>This is the party page for party {groupId}</>;
}
