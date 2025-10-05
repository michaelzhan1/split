import { useEffect, useRef, useState } from 'react';

import type { MultiSelectProps } from 'src/types/component.type';

import 'src/components/common/multiselect.component.css';

export function MultiSelect({
  id,
  options,
  selectedValues,
  onChange,
}: MultiSelectProps) {
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

  function toggleValue(value: string) {
    if (selectedValues.includes(value)) {
      onChange(selectedValues.filter((v) => v !== value));
    } else {
      onChange([...selectedValues, value]);
    }
  }

  return (
    <div ref={containerRef} className='multiselect-container'>
      <input
        id={id}
        className='multiselect-input'
        onClick={() => setIsOpen(!isOpen)}
        placeholder='Select'
        value={
          selectedValues.length === 0
            ? ''
            : options
                .filter((option) => selectedValues.includes(option.value))
                .map((option) => option.label)
                .join(', ')
        }
        readOnly
      />
      <div className='multiselect-list' hidden={!isOpen}>
        {options.map((option) => (
          <div
            key={option.value}
            className={`multiselect-item ${
              selectedValues.includes(option.value) ? 'selected' : ''
            }`}
            onClick={() => toggleValue(option.value)}
          >
            {option.label}
            {selectedValues.includes(option.value) && (
              <span className='checkmark'>✓</span>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
