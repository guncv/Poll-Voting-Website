import { CSSProperties, ChangeEvent } from 'react';

interface InputFieldProps {
  id: string;
  label: string;
  type: string;
  value: string;
  onChange: (event: ChangeEvent<HTMLInputElement>) => void;
}

export function InputField({ id, label, type, value, onChange }: InputFieldProps) {
    return (
      <div style={{ display: 'flex', flexDirection: 'column', gap: '15px' }}>
        <label htmlFor={id} style={{ color: 'black', fontSize: '20px' }}>{label}</label>
        <input
          id={id}
          type={type}
          placeholder={`Enter your ${label.toLowerCase()}`}
          value={value}
          onChange={onChange}
          style={inputStyle}
        />
      </div>
    );
}

const inputStyle: CSSProperties = {
  padding: '10px',
  borderRadius: '5px',
  border: '1px solid black',
  backgroundColor: 'white',
  color: 'black',
  fontSize: '15px',
  fontFamily: 'Poppins, sans-serif',
}