import axios, { type AxiosResponse } from 'axios';

import type { User } from 'src/types/common.type';

export async function getUsersByGroupId(groupId: number): Promise<User[]> {
  return axios
    .get<User[]>(`${import.meta.env.VITE_API_PREFIX}/groups/${groupId}/users`)
    .then((res) => res.data);
}

export async function addUserToGroup(
  groupId: number,
  name: string,
): Promise<{ id: number }> {
  return axios
    .post<
      { id: number },
      AxiosResponse,
      { name: string }
    >(`${import.meta.env.VITE_API_PREFIX}/groups/${groupId}/users`, { name })
    .then((res) => res.data);
}

export async function patchUser(
  groupId: number,
  id: number,
  name: string,
): Promise<void> {
  await axios.patch<void, AxiosResponse, { name: string }>(
    `${import.meta.env.VITE_API_PREFIX}/groups/${groupId}/users/${id}`,
    { name },
  );
}

export async function deleteUser(groupId: number, id: number): Promise<void> {
  await axios.delete<void, AxiosResponse, { name: string }>(
    `${import.meta.env.VITE_API_PREFIX}/groups/${groupId}/users/${id}`,
  );
}
