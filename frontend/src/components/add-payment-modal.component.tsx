import { useState } from 'react';
import { Dropdown } from 'src/components/common/dropdown.component';

import { Modal } from 'src/components/common/modal.component';
import type { AddPaymentModalProps } from 'src/types/component.type';

export function AddPaymentModal({
  isOpen,
  onClose,
  onSubmit,
  users,
}: AddPaymentModalProps) {
  const [description, setDescription] = useState<string>('');
  const [amountStr, setAmountStr] = useState<string>('');
  const [payerId, setPlayerId] = useState<number | null>(null);
  const [payeeIds, setPayeeIds] = useState<number[]>([]);

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title='Add Payment'
      onSubmit={() => {
        if (description.trim() === '') {
          alert('Description cannot be empty');
          return;
        }
        if (payerId === null) {
          alert('Payer must be selected');
          return;
        }
        if (payeeIds.length === 0) {
          alert('At least one payee must be selected');
          return;
        }
        if (isNaN(Number(amountStr)) || Number(amountStr) <= 0) {
          alert('Invalid amount');
          return;
        }
        onSubmit({
          description,
          amount: Number(amountStr),
          payer_id: payerId,
          payee_ids: payeeIds,
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
        <label htmlFor='payment-amountStr-input'>Amount</label>
        <input
          id='payment-amountStr-input'
          value={amountStr ?? ''}
          onChange={(e) => setAmountStr(e.target.value.replace(/\D/g, ''))}
        />
        <label htmlFor='payment-payer-id-input'>Payer ID</label>
        <Dropdown
          options={users.map((user) => ({ label: user.name, value: user.id.toString() }))}
          selectedValue={payerId?.toString() ?? ''}
          onSelect={(value) => setPlayerId(Number(value))}
        />
        {/* <input
          id='payment-payer-id-input'
          type='number'
          value={payerId ?? ''}
          onChange={(e) => setPlayerId(Number(e.target.value))}
        /> */}
        <label htmlFor='payment-payee-ids-input'>
          Payee IDs (comma separated)
        </label>
        <input
          id='payment-payee-ids-input'
          value={payeeIds.join(',')}
          onChange={(e) =>
            setPayeeIds(
              e.target.value
                .split(',')
                .map((id) => Number(id.trim()))
                .filter((id) => !isNaN(id)),
            )
          }
        />
      </form>
    </Modal>
  );
}
