import { createPortal } from 'react-dom';
import type { ModalProps } from 'src/types/component.type';

export function Modal({ isOpen, onClose, title, primaryActions }: ModalProps) {
  return createPortal(
    isOpen ? (
      <div className='modal-backdrop' onClick={onClose}>
        <div className='modal' onClick={(e) => e.stopPropagation()}>
          <div className='modal-header'>
            <h2>{title}</h2>
            <button className='modal-close-button' onClick={onClose}>
              &times;
            </button>
          </div>
          <div className='modal-body'>
            {/* Modal content can go here */}
            <p>This is a modal dialog.</p>
          </div>
          {primaryActions && primaryActions.length > 0 && (
            <div className='modal-footer'>
              {primaryActions.map((actionProps, index) => (
                <button key={index} {...actionProps} />
              ))}
            </div>
          )}
        </div>
      </div>
    ) : null,
    document.body,
  );
}
