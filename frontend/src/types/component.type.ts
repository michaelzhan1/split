import type { ComponentProps } from 'react';

export interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  primaryActions?: ComponentProps<'button'>[];
}
