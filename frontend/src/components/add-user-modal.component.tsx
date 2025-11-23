import { useState } from 'react';

import { Modal } from 'src/components/common/modal.component';
import type { AddUserModalProps } from 'src/types/component.type';

export function AddUserModal({ isOpen, onClose, onSubmit }: AddUserModalProps) {
  const [name, setName] = useState<string>('');

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title='Add User'
      onSubmit={() => {
        if (name === '') {
          alert('Name cannot be empty');
          return;
        }
        onSubmit(name);
      }}
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
