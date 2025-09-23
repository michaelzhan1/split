import { useState } from 'react';
import { useNavigate } from 'react-router';

export function Home() {
  const navigate = useNavigate();

  const [groupId, setGroupId] = useState<string>('');

  const onFindGroup = (e: React.FormEvent) => {
    e.preventDefault();

    navigate(`/groups/${groupId}`);
  };

  return (
    <>
      <h1>Split</h1>
      <h3>Find a group</h3>
      <form>
        <label htmlFor='group-code'>Group Code</label>
        <input
          name='group-code'
          type='number'
          placeholder='Group Code'
          value={groupId}
          onChange={(e) => setGroupId(e.target.value)}
        />
        <button type='submit' onClick={onFindGroup}>
          Go
        </button>
      </form>

      <h3>Create a new group</h3>
      <form>
        <label htmlFor='group-name'>Group Name</label>
        <input name='group-name' type='text' placeholder='Group Name' />
        <button type='submit'>Create</button>
      </form>
    </>
  );
}
