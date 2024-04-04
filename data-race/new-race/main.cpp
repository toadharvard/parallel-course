#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <vector>
#include <numeric>
#include <limits.h>
#include <cmath>
#include <iostream>
#include "fineGrainedSync.h"
#include "lazySync.h"

#define WRITERS_NUMB 3
#define READERS_NUMB 3

using namespace std;

int realSize = 2000;
vector<int> testIntArray(2*realSize, 0);

struct Args {
	int dataSize;
	vector<int> data;
	mySet<int> * set;
};

void *AddElem(void * args)
{
	Args *arguments = (Args*)args;

	for (int i = 0; i < arguments->dataSize; i++)
	{
		arguments->set->Add(arguments->data[i]);
	}
	pthread_exit(NULL);
	return NULL;
}


void *DeleteElem(void * args)
{
	Args* arguments = (Args*)args;
	for (int i = 0; i < arguments->dataSize; i++)
	{
		arguments->set->Remove(arguments->data[i]);
	}
	pthread_exit(NULL);
	return NULL;
}


void *ReadElem(void * args)
{
	Args* arguments = (Args*)args;
	for (int i = 0; i < arguments->dataSize; i++)
	{
		if(arguments->set->Contains(arguments->data[i]))
            		testIntArray.at(arguments->data[i])++;

	}
	pthread_exit(NULL);
	return NULL;
}


bool runTest(mySet<int> * set, int readersNum, int writersNum)
{
	bool writersTest = true;
	vector <pthread_t> writers(writersNum);
	vector <pthread_t> readers(readersNum);
	vector<int> addedData;

	for (int i = 0; i < 300; i++)
	{
		addedData.push_back(i);
	}

	int curpos = 0;
	cout << "starting writers test\n";
	for (int i = 0; i < writersNum; i++)
	{
		vector<int> dataToAdd;
		int dataSizeToAdd = (int)ceil((double)addedData.size() / (double)writersNum);
		if (addedData.size() <= (unsigned int)curpos)
			break;
		for (int j = 0; j < dataSizeToAdd; j++)
		{
			dataToAdd.push_back(addedData[curpos]);
			curpos ++;
		}

		Args *argsStruct;
		argsStruct = new Args;

		argsStruct->dataSize = dataToAdd.size();
		argsStruct->data = dataToAdd;
		argsStruct->set = set;
		
		pthread_create(&writers[i], NULL, AddElem, (void *)argsStruct);
		dataToAdd.resize(0);
	}

	for (int i = 0; i < writersNum; i++)
	{
		pthread_join(writers[i], NULL);
	}

	for (unsigned int i = 0; i < addedData.size(); i++)
	{
		if (!set->Contains(addedData[i])) {
			writersTest = false;
			break;
		}
	}

	if (writersTest)
		cout << "writers test: ok\n";
	else
		cout << "writers test: failed\n";

	cout << "starting readers writers test\n";


    	bool readersWritersTest = true;

	for (int i = 0; i < realSize; i ++)
		set->Add(2*i);

    	for (int i = 0; i < writersNum; i++)
	{
		vector<int> dataToAdd2;
		addedData.push_back(2 * i + 1);

		Args *argsStruct;
		argsStruct = new Args;

		argsStruct->dataSize = dataToAdd2.size();
		argsStruct->data = dataToAdd2;
		argsStruct->set = set;

		pthread_create(&writers[i], NULL, AddElem, (void *)argsStruct);
		dataToAdd2.resize(0);
	}

	curpos = 0;
	int runed = 0;
	for (int i = 0; i < readersNum; i++)
	{
		vector<int> dataToRemove;
		int dataSizeToRemove = (int)ceil((float)realSize / (float)readersNum);
		if (realSize <= curpos)
			break;
		for (int j = 0; j < dataSizeToRemove; j++)
		{
			dataToRemove.push_back(2*curpos);
			curpos ++;
		}
		Args *argsStruct;
		argsStruct = new Args;

		argsStruct->dataSize = dataToRemove.size();
		argsStruct->data = dataToRemove;
		argsStruct->set = set;

		pthread_create(&readers[i], NULL, ReadElem, (void *)argsStruct);
		runed++;
		dataToRemove.resize(0);
	}

	for (int i = 0; i < runed; i++)
	{
		pthread_join(readers[i], NULL);
	}

	for (int i = 0; i < writersNum; i++)
	{
		pthread_join(writers[i], NULL);
	}

	for (int i = 0; i < realSize; i ++)
	{
        	if(testIntArray[2*i] != 1)
        	{
            		readersWritersTest = false;
        	}
	}

	if (readersWritersTest)
		cout << "readers writers test: ok\n";
	else
		cout << "readers writers test: failed\n";


	cout << "starting readers test\n";


	return true;
}

int main (int argc, char * argv [])
{
	mySet<int> *set;


    //     set = new FGSet<int>(INT_MIN, INT_MAX);
	// cout << "TEST FOR FINE GRAINDES SET\n";
	// runTest(set, READERS_NUMB, WRITERS_NUMB);
	// cout << "TEST END\n\n";

	// delete set;

	testIntArray.resize(0);
	testIntArray.resize(2*realSize, 0);

        set = new LSSet<int>(INT_MIN, INT_MAX);

	cout << "TEST FOR LAZY SET\n";
	runTest(set, READERS_NUMB, WRITERS_NUMB);
	cout << "TEST END\n\n";

	delete set;
	return 0;
}
