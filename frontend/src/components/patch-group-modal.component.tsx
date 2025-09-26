import { useState } from 'react';

import { Modal } from 'src/components/common/modal.component';
import type { PatchGroupModalProps } from 'src/types/component.type';

export function PatchGroupModal({
  isOpen,
  onClose,
  onSubmit,
  initialName,
}: PatchGroupModalProps) {
  const [name, setName] = useState<string>(initialName);

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title='Edit Group'
      onSubmit={() => onSubmit(name)}
    >
      <form>
        <label htmlFor='group-name-input'>Name</label>
        <input
          id='group-name-input'
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
      </form>
    </Modal>
  );
}
