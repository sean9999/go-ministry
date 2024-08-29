import { uuidv7 } from "uuidv7";

export type SoccerMessageHandler = (m: SoccerMessage) => void;

export type SoccerRecord = {
    id: string;
    thread_id?: string;
    subject: string;
    payload: any;
};

export class SoccerMessage{
    public record : SoccerRecord;
    constructor(subj="", pay=null){
        this.record = {
            id: uuidv7(),
            thread_id: null,
            subject: subj,
            payload: pay
        };
    }
    get subject() : string {
        return this.record.subject;
    }
    get thread_id() : string | null | undefined {
        return this.record.thread_id;
    }
    get id() : string {
        return this.record.id;
    }
    get payload() : any {
        return this.record.payload;
    }
    static deserialize(txt : string) : SoccerMessage {
        const rec = JSON.parse(txt)
        const msg = new SoccerMessage();
        msg.record = rec;
        return msg;
    }
    serialize() : string {
        return JSON.stringify(this.record);
    }
    reply(pay=null) : SoccerMessage {
        const s = new SoccerMessage(this.record.subject, pay);
        if (this.record.thread_id) {
            s.record.thread_id = this.record.thread_id;
        } else {
            s.record.thread_id = this.record.id;
        }
        return s;
    }
}

