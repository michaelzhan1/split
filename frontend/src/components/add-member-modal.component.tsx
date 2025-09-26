import { useState } from 'react';

import { Modal } from 'src/components/common/modal.component';
import type { AddMemberModalProps } from 'src/types/component.type';

export function AddMemberModal({
  isOpen,
  onClose,
  onSubmit,
}: AddMemberModalProps) {
  const [name, setName] = useState<string>('');

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title='Add Member'
      onSubmit={() => onSubmit(name)}
    >
      <form>
        <label htmlFor="member-name-input">Name</label>
        <input id="member-name-input" value={name} onChange={(e) => setName(e.target.value)} />
      </form>
    </Modal>
  );
}
