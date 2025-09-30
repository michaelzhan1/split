import { Modal } from 'src/components/common/modal.component';
import type { ConfirmationModalProps } from 'src/types/component.type';

export function ConfirmationModal({
  isOpen,
  onClose,
  title,
  content,
  onSubmit,
}: ConfirmationModalProps) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} title={title} onSubmit={onSubmit}>
      {content}
    </Modal>
  );
}
