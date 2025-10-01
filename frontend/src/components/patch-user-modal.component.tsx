import { useState } from 'react';

import { Modal } from 'src/components/common/modal.component';
import type { PatchUserModalProps } from 'src/types/component.type';

export function PatchUserModal({
  isOpen,
  onClose,
  onSubmit,
  user,
}: PatchUserModalProps) {
  const [name, setName] = useState<string>(user.name);

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title='Edit User'
      onSubmit={() => onSubmit(user.id, name)}
    >
      <form>
        <label htmlFor='user-name-input'>Name</label>
        <input
          id='user-name-input'
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
      </form>
    </Modal>
  );
}
