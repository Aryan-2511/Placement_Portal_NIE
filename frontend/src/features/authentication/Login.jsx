// import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useState } from 'react';
import toast from 'react-hot-toast';

import { Input } from '../../components/ui/input';
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from '../../components/ui/card';
import { Label } from '../../components/ui/label';
import { Button } from '../../components/ui/button';
import { RadioGroup, RadioGroupItem } from '../../components/ui/radio-group';
import { verifyLogin } from '@/utils/verifyLogin';
import useLogin from './useLogin';
import Spinner from '@/components/shared/Spinner';

function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('student');
  const { login, isLoading } = useLogin();
  const [error, setError] = useState('');

  function handleSubmit() {
    const result = verifyLogin(email, password);
    if (!result.valid) {
      setError(result.error);
      toast.error(result.error);
    } else {
      login(
        { email, password, role },
        {
          onSettled: () => {
            setError('');
            setEmail('');
            setPassword('');
            setRole('STUDENT');
          },
        }
      );
    }
  }

  if (isLoading) return <Spinner />;

  return (
    <Card className="px-[2.4rem] py-[1.2rem] text-[var(--color-grey-600)]">
      <CardHeader>
        <CardTitle className="text-[2rem]">Login to Your Account</CardTitle>
        {/* <CardDescription className="text-[1.6rem]">
          Login with email and password. Register if account doesn&apos;t
          already exists.
        </CardDescription> */}
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="email" className="">
            Email
          </Label>
          <Input
            id="email"
            type="email"
            placeholder="johndoe@email.com"
            onChange={(e) => setEmail(e.target.value)}
            value={email}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="password">Password</Label>
          <Input
            id="password"
            type="password"
            placeholder="password@1234"
            onChange={(e) => setPassword(e.target.value)}
            value={password}
          />
        </div>
        <RadioGroup
          defaultValue={role}
          value={role}
          className="space-y-2"
          onValueChange={(value) => setRole(value)}
        >
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="STUDENT" id="student" />
            <Label htmlFor="student">Student</Label>
          </div>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="COORDINATOR" id="coordinator" />
            <Label htmlFor="coordinator">Coordinator</Label>
          </div>
          <div className="flex items-center space-x-2">
            <RadioGroupItem value="ADMIN" id="admin" />
            <Label htmlFor="admin">Admin</Label>
          </div>
        </RadioGroup>
      </CardContent>
      <CardFooter className="flex mt-4">
        <Button className="w-1/2" variant="ghost">
          REGISTER
        </Button>
        <Button className="w-1/2" onClick={() => handleSubmit()}>
          LOGIN
        </Button>
      </CardFooter>
    </Card>
  );
}

export default Login;
