import React from 'react';
import {
    Anchor,
    Box,
    Button,
    Center,
    Combobox,
    Container,
    Group,
    Paper,
    PasswordInput,
    Progress,
    SegmentedControl,
    Space,
    Text,
    TextInput,
    Title
} from '@mantine/core';
import {getHotkeyHandler, useInputState} from "@mantine/hooks";
import {IconAlertTriangle, IconCheck, IconX} from '@tabler/icons-react';
import {IUsersAdapter} from "../../Adapters/IUsersAdapter";
import {User} from "../../Models/User";
import classes = Combobox.classes;
import {useNavigate} from "react-router-dom";

interface IProps {
    usersAdapter: IUsersAdapter;
}

function PasswordRequirement({meets, label}: { meets: boolean; label: string }) {
    return (
        <Text component="div" c={meets ? 'teal' : 'red'} mt={5} size="sm">
            <Center inline>
                {meets ? <IconCheck size={14} stroke={1.5}/> : <IconX size={14} stroke={1.5}/>}
                <Box ml={7}>{label}</Box>
            </Center>
        </Text>
    );
}

const requirements = [
    {re: /[0-9]/, label: 'Includes number'},
    {re: /[a-z]/, label: 'Includes lowercase letter'},
    {re: /[A-Z]/, label: 'Includes uppercase letter'},
];

function getStrength(password: string) {
    let multiplier = password.length > 5 ? 0 : 1;

    requirements.forEach((requirement) => {
        if (!requirement.re.test(password)) {
            multiplier += 1;
        }
    });

    return Math.max(100 - (100 / (requirements.length + 1)) * multiplier, 0);
}

export const LoginCard: React.FunctionComponent<IProps> = (props) => {
    const [registrationPassword, setRegistrationPassword] = useInputState('');
    const strength = getStrength(registrationPassword);
    const [segmentedValue, setSegmentedValue] = useInputState('Login');
    const [username, setUsername] = useInputState('');
    const [email, setEmail] = useInputState('');
    const [password, setPassword] = useInputState('');
    const navigate = useNavigate()
    const checks = requirements.map((requirement, index) => (
        <PasswordRequirement key={index} label={requirement.label} meets={requirement.re.test(registrationPassword)}/>
    ));
    const bars = Array(4)
        .fill(0)
        .map((_, index) => (
            <Progress
                styles={{section: {transitionDuration: '0ms'}}}
                value={
                    registrationPassword.length > 0 && index === 0 ? 100 : strength >= ((index + 1) / 4) * 100 ? 100 : 0
                }
                color={strength > 80 ? 'teal' : strength > 50 ? 'yellow' : 'red'}
                key={index}
                size={4}
            />
        ));

    const handleSubmit = async () => {
        let token;
        switch (segmentedValue) {
            case 'Login':
                token = await props.usersAdapter.login(new User(username, "", password));
                break;
            case 'Register':
                token = await props.usersAdapter.register(new User(username, email, registrationPassword));
                break;
        }
        localStorage.setItem("token", token.token);
        navigate("/")
    };

    const loginForm = () => {
        return (
            <>
                <TextInput
                    value={username}
                    label="Username"
                    error="Invalid username"
                    placeholder="username"
                    width="75%"
                    required
                    onChange={setUsername}
                    rightSection={<IconAlertTriangle stroke={1.5} size={18} className={classes.icon}/>}
                />
                <PasswordInput
                    value={password}
                    placeholder="Your password"
                    label="Password"
                    required
                    onChange={setPassword}
                    width="75%"
                    onKeyDown={getHotkeyHandler([
                        ['Enter', handleSubmit],
                    ])}
                />
            </>
        );
    }

    const registrationForm = () => {
        return (
            <>
                <TextInput
                    value={username}
                    label="Username"
                    error="Invalid username"
                    placeholder="username"
                    width="75%"
                    required
                    onChange={setUsername}
                    rightSection={<IconAlertTriangle stroke={1.5} size={18} className={classes.icon}/>}
                />
                <TextInput
                    value={email}
                    label="Email"
                    error="Invalid email"
                    placeholder="hello@gmail.com"
                    width="75%"
                    onChange={setEmail}
                    rightSection={<IconAlertTriangle stroke={1.5} size={18} className={classes.icon}/>}
                />
                <PasswordInput
                    value={registrationPassword}
                    onChange={setRegistrationPassword}
                    placeholder="Your password"
                    label="Password"
                    required
                    width="75%"
                    onKeyDown={getHotkeyHandler([
                        ['Enter', handleSubmit],
                    ])}
                />
                <Group gap={5} grow mt="xs" mb="md">
                    {bars}
                </Group>

                <PasswordRequirement label="Has at least 6 characters" meets={registrationPassword.length > 5}/>
                {checks}
            </>
        )
    }

    return (
        <Container size={420} my={40}>
            <Title ta="center" className={classes.title}>
                Welcome back!
            </Title>
            <Text c="dimmed" size="sm" ta="center" mt={5}>
                Do not have an account yet?{' '}
                <Anchor size="sm" component="button">
                    Create account
                </Anchor>
            </Text>
            <Paper withBorder shadow="md" p={30} mt={30} radius="md">
                <SegmentedControl data={['Login', 'Register']} fullWidth value={segmentedValue} onChange={setSegmentedValue}/>
                <Space h={20}/>
                {segmentedValue === 'Login' ? loginForm() : registrationForm()}
                <Group justify="space-between" mt="lg">
                    <Anchor component="button" size="sm">
                        Forgot password?
                    </Anchor>
                </Group>
                <Button fullWidth mt="xl" onClick={() => handleSubmit()}>
                    {segmentedValue}
                </Button>
            </Paper>
        </Container>
    );
};
