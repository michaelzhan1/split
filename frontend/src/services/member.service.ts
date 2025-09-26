import axios, { type AxiosResponse } from 'axios';

import type { Member } from 'src/types/common.type';

export async function getMembersByGroupId(id: number): Promise<Member[]> {
  return axios
    .get<Member[]>(`${import.meta.env.VITE_API_PREFIX}/parties/${id}/members`)
    .then((res) => res.data);
}

export async function addMembertoGroup(
  id: number,
  name: string,
): Promise<{ id: number }> {
  return axios
    .post<
      { id: number },
      AxiosResponse,
      { name: string }
    >(`${import.meta.env.VITE_API_PREFIX}/parties/${id}/members`, { name })
    .then((res) => res.data);
}
