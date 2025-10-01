import axios, { type AxiosResponse } from 'axios';

import type { Member } from 'src/types/common.type';

export async function getUsersByGroupId(id: number): Promise<Member[]> {
  return axios
    .get<Member[]>(`${import.meta.env.VITE_API_PREFIX}/groups/${id}/users`)
    .then((res) => res.data);
}

export async function addUserToGroup(
  id: number,
  name: string,
): Promise<{ id: number }> {
  return axios
    .post<
      { id: number },
      AxiosResponse,
      { name: string }
    >(`${import.meta.env.VITE_API_PREFIX}/groups/${id}/users`, { name })
    .then((res) => res.data);
}

export async function patchUser(
  partyId: number,
  id: number,
  name: string,
): Promise<void> {
  await axios.patch<void, AxiosResponse, { name: string }>(
    `${import.meta.env.VITE_API_PREFIX}/groups/${partyId}/users/${id}`,
    {
      name,
    },
  );
}

export async function deleteUser(partyId: number, id: number): Promise<void> {
  await axios.delete<void, AxiosResponse, { name: string }>(
    `${import.meta.env.VITE_API_PREFIX}/groups/${partyId}/users/${id}`,
  );
}
