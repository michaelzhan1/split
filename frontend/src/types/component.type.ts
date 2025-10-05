import type { ReactNode } from 'react';
import type { CreatePaymentRequest, User } from 'src/types/common.type';

interface DropdownOption {
  label: string;
  value: string;
}

export interface DropdownProps {
  options: DropdownOption[];
  selectedValue: string;
  onSelect: (value: string) => void;
}

export interface MultiSelectProps {
  options: DropdownOption[];
  selectedValues: string[];
  onChange: (values: string[]) => void;
}

export interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  onSubmit: () => void;
  children: ReactNode;
}

export interface AddUserModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (name: string) => void;
}

export interface PatchGroupModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (name: string) => void;
  initialName: string;
}

export interface PatchUserModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (id: number, name: string) => void;
  user: User;
}

export interface AddPaymentModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (data: CreatePaymentRequest) => void;
  users: User[];
}

export interface ConfirmationModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  content: string;
  onSubmit: () => void;
}
