import type { ReactNode } from 'react';
import type { Member } from 'src/types/common.type';

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

export interface PatchGroupModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (name: string) => void;
  initialName: string;
}

export interface PatchMemberModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (id: number, name: string) => void;
  member: Member;
}
