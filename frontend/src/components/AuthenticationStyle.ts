import { CSSProperties } from 'react';

export const formStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  gap: '20px',
  width: '500px',
};

export const cardStyle: CSSProperties = {
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'space-between',
  border: '1px solid #FF98A0',
  padding: '40px',
  borderRadius: '10px',
};

export const titleStyle: CSSProperties = {
  fontSize: '50px',
  fontWeight: 'bold',
  color: 'black',
  marginBottom: '50px',
};

export const inputStyle: CSSProperties = {
  padding: '10px',
  borderRadius: '5px',
  border: '1px solid black',
  backgroundColor: 'white',
  color: 'black',
  fontSize: '15px',
};

export const buttonStyle: CSSProperties = {
  marginTop: '20px',
  backgroundColor: '#FF98A0',
  color: 'white',
  padding: '10px',
  borderRadius: '10px',
  border: 'none',
  fontSize: '20px',
  fontFamily: 'Poppins, sans-serif',
};

export const linkStyle: CSSProperties = {
  fontSize: '14px',
  cursor: 'pointer',
};
