declare global {
    interface Window {
        Telegram: {
            WebApp: WebApp;
        };
    }
}

interface WebApp {
    initData: any;
    expand(): void;
    enableClosingConfirmation(): void;
    MainButton: {
        setParams(params: { color: string; text: string }): void;
        show(): void;
    };
    HapticFeedback: {
        notificationOccurred(type: string): void;
        selectionChanged(): void;
        impactOccurred(light: string): void;
    };

    onEvent(mainButtonClicked: string, save: () => Promise<void>): void;

    disableVerticalSwipes(): void;

    close(): void;
}

class Telegram {
    webapp: WebApp;
    constructor() {
        this.webapp = window.Telegram.WebApp;
    }

    initApp(): void {
        this.webapp.expand();
        this.webapp.enableClosingConfirmation();

    }

    showSaveButton(): void {
        this.webapp.MainButton.setParams({
            color: '#242424',
            text: 'Save'
        });
        this.webapp.MainButton.show();
    }
}

const telegram = new Telegram();
export default telegram;