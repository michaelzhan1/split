import { useQuery } from '@tanstack/react-query';

import type { AxiosError } from 'axios';
import { useParams } from 'react-router';
import { getGroupById } from 'src/services/group.service';
import type { Group } from 'src/types/common.type';

export function Group() {
  const { groupId = '' } = useParams();

  const { data: group = null, isFetching: isFetchingGroup } = useQuery<
    Group,
    AxiosError
  >({
    queryKey: ['group', groupId],
    queryFn: () => getGroupById(Number(groupId)),
  });

  return (
    !group || isFetchingGroup ? (
      <div>Loading...</div>
    ) : (
      <div>
        <h1>Group: {group.name}</h1>
      </div>
    )
  );
}
