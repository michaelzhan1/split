import { createPortal } from 'react-dom';

import type { ModalProps } from 'src/types/component.type';

import 'src/components/modal.component.css';

export function Modal({
  isOpen,
  onClose,
  title,
  children,
}: ModalProps) {
  return createPortal(
    isOpen ? (
      <div className='modal-backdrop' onClick={onClose}>
        <div className='modal' onClick={(e) => e.stopPropagation()}>
          <div className='modal-header'>
            <h2>{title}</h2>
            <button onClick={onClose}>&times;</button>
          </div>
          <hr />
          <div className='modal-body'>{children}</div>
          <div className='modal-buttons'>
            <button className='submit-button'>Submit</button>
            <button className='cancel-button' onClick={onClose}>
              Cancel
            </button>
          </div>
        </div>
      </div>
    ) : null,
    document.body,
  );
}
