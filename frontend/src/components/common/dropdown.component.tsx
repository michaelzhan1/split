import { useEffect, useRef, useState } from 'react';

import type { DropdownProps } from 'src/types/component.type';

import 'src/components/common/dropdown.component.css';

export function Dropdown({
  id,
  options,
  selectedValue,
  onSelect,
}: DropdownProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [isOpen, setIsOpen] = useState(false);

  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        containerRef.current &&
        !containerRef.current.contains(event.target as HTMLElement)
      ) {
        setIsOpen(false);
      }
    }
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    <div ref={containerRef} className='dropdown-container'>
      <input
        id={id}
        className='dropdown-input'
        value={selectedValue}
        readOnly
        onClick={() => setIsOpen(!isOpen)}
      />
      <div className='dropdown-list' hidden={!isOpen}>
        {options.map((option) => (
          <div
            key={option.value}
            className='dropdown-item'
            onClick={() => {
              onSelect(option.value);
              setIsOpen(false);
            }}
          >
            {option.label}
          </div>
        ))}
      </div>
    </div>
  );
}
