import { SoccerMessage } from '../src/js/msg';

test("well formed message", () => {
    let msg = new SoccerMessage("test");
    expect(msg.subject).toEqual("test");  
    expect(msg.thread_id).toBeFalsy();
    expect(msg.payload).toBeNull();
});

test("well formed reply", () => {
    let msg = new SoccerMessage("test");
    let reply = msg.reply();
    expect(reply.subject).toEqual(msg.subject);
    expect(reply.thread_id).toEqual(msg.id);
});