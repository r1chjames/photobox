import React from 'react';

const MockIcon = React.forwardRef<SVGSVGElement, any>((props, ref) =>
    React.createElement('svg', { ref, 'data-testid': 'mock-icon', ...props })
);

const handler = {
    get(target: any, prop: string) {
        if (prop === 'default' || prop === '__esModule') {
            return target[prop];
        }
        return MockIcon;
    },
};

export default new Proxy({}, handler);

// Named exports via proxy
export const IconPhotoOff = MockIcon;
export const IconRefresh = MockIcon;
export const IconPlayerPlay = MockIcon;
export const IconPhoto = MockIcon;
export const IconAlbum = MockIcon;
export const IconHeart = MockIcon;
export const IconTrash = MockIcon;
export const IconLink = MockIcon;
export const IconUsers = MockIcon;
export const IconMap = MockIcon;
export const IconTag = MockIcon;
export const IconCopy = MockIcon;
export const IconVideo = MockIcon;
export const IconSettings = MockIcon;
export const IconHome = MockIcon;
export const IconLibraryPhoto = MockIcon;
export const IconSearch = MockIcon;
export const IconChevronDown = MockIcon;
export const IconMoon = MockIcon;
export const IconSun = MockIcon;
export const IconLogout = MockIcon;
export const IconMessage = MockIcon;
export const IconHelp = MockIcon;
export const IconLayoutGrid = MockIcon;
export const IconList = MockIcon;
export const IconSelect = MockIcon;
export const IconCheck = MockIcon;
export const IconWifiOff = MockIcon;
export const IconServerOff = MockIcon;
export const IconMenu2 = MockIcon;
export const IconArrowLeftDashed = MockIcon;
export const IconArrowRightDashed = MockIcon;
export const IconCalendar = MockIcon;
export const IconFolder = MockIcon;
export const IconDownload = MockIcon;
export const IconHeartFilled = MockIcon;
export const IconMaximize = MockIcon;
export const IconMinimize = MockIcon;
export const IconRotateClockwise = MockIcon;
export const IconShare2 = MockIcon;
export const IconX = MockIcon;
export const IconRestore = MockIcon;
export const IconPhotoPlus = MockIcon;
export const IconDots = MockIcon;
export const IconClock = MockIcon;
export const IconInfoCircle = MockIcon;
export const IconHome2 = MockIcon;
