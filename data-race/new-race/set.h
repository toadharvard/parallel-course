#ifndef SET_H
#define SET_H

#include <pthread.h>

template <typename Elem>
class mySet {
public:
   mySet ()  {}
   virtual ~mySet()  {}

   virtual bool  Add       (const Elem & item) = 0;
   virtual bool  Remove    (const Elem & item) = 0;
   virtual bool  Contains  (const Elem & item) = 0;
};

template <typename Elem>
class Node {
public:
	Node(Elem val)
	{
		data = val;
		deleted = false;
		next = NULL;
		pthread_mutex_init(&mutex, NULL);
	}
	~Node(void)
	{
		pthread_mutex_destroy(&mutex);
	}

	Elem data;
	Node *next;
	pthread_mutex_t mutex;
	bool deleted;
};

#endif
