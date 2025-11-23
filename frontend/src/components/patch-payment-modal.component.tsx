import { useState } from 'react';

import { Modal } from 'src/components/common/modal.component';
import type { PatchPaymentModalProps } from 'src/types/component.type';

export function PatchPaymentModal({
  isOpen,
  onClose,
  onSubmit,
  payment,
}: PatchPaymentModalProps) {
  const [description, setDescription] = useState<string>(payment.description);
  const [amountStr, setAmountStr] = useState<string>(payment.amount.toString());

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title='Edit Payment'
      onSubmit={() => {
        if (description.trim() === '') {
          alert('Description cannot be empty');
          return;
        }
        if (isNaN(Number(amountStr)) || Number(amountStr) <= 0) {
          alert('Invalid amount');
          return;
        }
        onSubmit(payment.id, {
          id: payment.id,
          description,
          amount: Number(amountStr),
        });
      }}
    >
      <form>
        <label htmlFor='payment-description-input'>Description</label>
        <input
          id='payment-description-input'
          value={description}
          onChange={(e) => setDescription(e.target.value)}
        />
        <label htmlFor='payment-amount-input'>Amount</label>
        <input
          id='payment-amount-input'
          value={amountStr}
          onChange={(e) => setAmountStr(e.target.value.replace(/\D/g, ''))}
        />
      </form>
    </Modal>
  );
}
