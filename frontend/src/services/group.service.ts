import axios from 'axios';
import type { Group } from 'src/types/common.type';

export async function getGroupById(id: number): Promise<Group> {
  return axios
    .get<Group>(`${import.meta.env.VITE_API_PREFIX}/parties/${id}`)
    .then((res) => res.data);
}
