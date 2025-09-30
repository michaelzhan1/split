import { useState } from 'react';

import { Modal } from 'src/components/common/modal.component';
import type { PatchMemberModalProps } from 'src/types/component.type';

export function PatchMemberModal({
  isOpen,
  onClose,
  onSubmit,
  member,
}: PatchMemberModalProps) {
  const [name, setName] = useState<string>(member.name);

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title='Edit Member'
      onSubmit={() => onSubmit(member.id, name)}
    >
      <form>
        <label htmlFor='member-name-input'>Name</label>
        <input
          id='member-name-input'
          value={name}
          onChange={(e) => setName(e.target.value)}
        />
      </form>
    </Modal>
  );
}
