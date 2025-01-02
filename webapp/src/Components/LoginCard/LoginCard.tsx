import React from 'react';
import {
    Box,
    Button,
    Card,
    Center,
    Combobox,
    Group,
    Modal,
    PasswordInput,
    Progress, SegmentedControl,
    Text,
    TextInput
} from '@mantine/core';
import {useInputState} from "@mantine/hooks";
import {IconAlertTriangle, IconCheck, IconX} from '@tabler/icons-react';
import classes = Combobox.classes;

interface IProps {
    isLoggedIn: boolean;
    submitLogin: () => void;
}

function PasswordRequirement({ meets, label }: { meets: boolean; label: string }) {
    return (
        <Text component="div" c={meets ? 'teal' : 'red'} mt={5} size="sm">
            <Center inline>
                {meets ? <IconCheck size={14} stroke={1.5} /> : <IconX size={14} stroke={1.5} />}
                <Box ml={7}>{label}</Box>
            </Center>
        </Text>
    );
}

const requirements = [
    { re: /[0-9]/, label: 'Includes number' },
    { re: /[a-z]/, label: 'Includes lowercase letter' },
    { re: /[A-Z]/, label: 'Includes uppercase letter' },
    { re: /[$&+,:;=?@#|'<>.^*()%!-]/, label: 'Includes special symbol' },
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
    const [value, setValue] = useInputState('');
    const [password, setPassword] = useInputState('');
    const [username, setUsername] = useInputState('');
    const [email, setEmail] = useInputState('');
    const strength = getStrength(value);
    const [segmentedValue, setSegmentedValue] = useInputState('Login');
    const checks = requirements.map((requirement, index) => (
        <PasswordRequirement key={index} label={requirement.label} meets={requirement.re.test(value)} />
    ));
    const bars = Array(4)
        .fill(0)
        .map((_, index) => (
            <Progress
                styles={{ section: { transitionDuration: '0ms' } }}
                value={
                    value.length > 0 && index === 0 ? 100 : strength >= ((index + 1) / 4) * 100 ? 100 : 0
                }
                color={strength > 80 ? 'teal' : strength > 50 ? 'yellow' : 'red'}
                key={index}
                size={4}
            />
        ));

    const loginForm = () => {
        return (
            <>
                <TextInput
                    label="Username"
                    error="Invalid username"
                    defaultValue="username"
                    classNames={{ input: classes.invalid }}
                    width="75%"
                    required
                    rightSection={<IconAlertTriangle stroke={1.5} size={18} className={classes.icon} />}
                />
                <PasswordInput
                    value={value}
                    onChange={setValue}
                    placeholder="Your password"
                    label="Password"
                    required
                    width="75%"
                />
            </>
        );
    }

    const registrationForm = () => {
        return (
            <>
                <TextInput
                    label="Username"
                    error="Invalid username"
                    defaultValue="username"
                    classNames={{ input: classes.invalid }}
                    width="75%"
                    required
                    rightSection={<IconAlertTriangle stroke={1.5} size={18} className={classes.icon} />}
                />
                <TextInput
                    label="Email"
                    error="Invalid email"
                    defaultValue="hello@gmail.com"
                    classNames={{ input: classes.invalid }}
                    width="75%"
                    rightSection={<IconAlertTriangle stroke={1.5} size={18} className={classes.icon} />}
                />
                <PasswordInput
                    value={value}
                    onChange={setValue}
                    placeholder="Your password"
                    label="Password"
                    required
                    width="75%"
                />
                <Group gap={5} grow mt="xs" mb="md">
                    {bars}
                </Group>

                <PasswordRequirement label="Has at least 6 characters" meets={value.length > 5}/>
                {checks}
            </>
        )
    }

    return (
        <Modal opened={!props.isLoggedIn} onClose={props.submitLogin} withCloseButton={false} centered>
            <Card shadow="sm" radius="md" padding={"xs"} mih={"400"}>
                <SegmentedControl data={['Login', 'Register']} value={segmentedValue} onChange={(value: string) => setSegmentedValue(value)}/>
                <Card.Section w={"100%"}>
                    <div>
                        {segmentedValue === 'Login' ? loginForm() : registrationForm()}


                    </div>
                </Card.Section>
                 <Button variant="filled" onClick={() => props.submitLogin}>Submit</Button>
            </Card>
        </Modal>
    );
};
