import NfcManager, { NfcTech } from 'react-native-nfc-manager';

class NFCService {
  static async init() {
    await NfcManager.start();
  }

  static async readNFC(): Promise<string> {
    try {
      await NfcManager.requestTechnology(NfcTech.Ndef);
      const tag = await NfcManager.getTag();
      return tag?.id ?? '';
    } catch (error) {
      console.error('Error reading NFC:', error);
      throw error;
    } finally {
      NfcManager.cancelTechnologyRequest();
    }
  }
}

export default NFCService;