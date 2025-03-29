import { ReactNode } from 'react';
import { cardStyle, titleStyle } from './AuthenticationStyle';

type FormLayoutProps = {
  title: string;
  children: ReactNode;
};

export default function FormLayout({ title, children }: FormLayoutProps) {
  return (
    <div style={cardStyle}>
      <h1 style={titleStyle}>{title}</h1>
      {children}
    </div>
  );
}
