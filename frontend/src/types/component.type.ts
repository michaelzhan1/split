import type { ReactNode } from 'react';

export interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  onSubmit: () => void;
  children: ReactNode;
}

export interface AddMemberModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (name: string) => void;
}
