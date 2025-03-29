import RegisterForm from '@/components/RegisterForm';

export default function RegisterPage() {
  return (
    <div style={{
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      height: '100vh'
    }}>
      <RegisterForm />
    </div>
  );
}
