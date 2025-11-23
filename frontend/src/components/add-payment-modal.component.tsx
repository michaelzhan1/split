import { useState } from 'react';

import { Dropdown } from 'src/components/common/dropdown.component';
import { Modal } from 'src/components/common/modal.component';
import { MultiSelect } from 'src/components/common/multiselect.component';
import type { AddPaymentModalProps } from 'src/types/component.type';

export function AddPaymentModal({
  isOpen,
  onClose,
  onSubmit,
  users,
}: AddPaymentModalProps) {
  const [description, setDescription] = useState<string>('');
  const [amountStr, setAmountStr] = useState<string>('');
  const [payerId, setPlayerId] = useState<string>('');
  const [payeeIds, setPayeeIds] = useState<string[]>([]);

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
        if (isNaN(Number(payerId)) || Number(payerId) <= 0) {
          alert('Invalid payer ID');
          return;
        }
        if (payeeIds.length === 0) {
          alert('At least one payee must be selected');
          return;
        }
        payeeIds.forEach((id) => {
          if (isNaN(Number(id)) || Number(id) <= 0) {
            alert(`Invalid payee ID: ${id}`);
            return;
          }
        });
        if (isNaN(Number(amountStr)) || Number(amountStr) <= 0) {
          alert('Invalid amount');
          return;
        }
        onSubmit({
          description,
          amount: Number(amountStr),
          payer_id: Number(payerId),
          payee_ids: payeeIds.map((id) => Number(id)),
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
        <label htmlFor='payment-payer-id-input'>Payer</label>
        <Dropdown
          id='payment-payer-id-input'
          options={users.map((user) => ({
            label: user.name,
            value: user.id.toString(),
          }))}
          selectedValue={payerId?.toString() ?? ''}
          onSelect={(value) => setPlayerId(value)}
        />
        <label htmlFor='payment-payee-ids-input'>
          Payees (comma separated)
        </label>
        <MultiSelect
          id='payment-payee-ids-input'
          options={users.map((user) => ({
            label: user.name,
            value: user.id.toString(),
          }))}
          selectedValues={payeeIds}
          onChange={(values: string[]) => setPayeeIds(values)}
        />
      </form>
    </Modal>
  );
}
