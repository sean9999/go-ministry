import { uuidv7 } from "uuidv7";

export type SoccerMessageHandler = (m: SoccerMessage) => void;

export type SoccerRecord = {
    from?: string;
    to?: string;
    id: string;
    thread_id?: string;
    subject: string;
    payload?: any;
};

export class SoccerMessage{
    public record : SoccerRecord;

    //  to do: simplify constructor to just make a blank object
    //  allow modification by simply setting properties.
    constructor(subj="", pay=null){
        this.record = {
            id: uuidv7(),
            thread_id: null,
            subject: subj,
            payload: pay
        };
    }

    //  to do: these accessors are not needed. It just makes everything to complicated
    get from() : string {
        return this.record.from;
    }
    get to() : string {
        return this.record.to;
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
        const retort = new SoccerMessage(this.record.subject, pay);
        if (this.from) {
            retort.record.to = this.from;
        }
        if (this.to) {
            retort.record.from = this.to;
        }
        if (this.thread_id) {
            retort.record.thread_id = this.record.thread_id;
        } else {
            retort.record.thread_id = this.record.id;
        }
        return retort;
    }
}

